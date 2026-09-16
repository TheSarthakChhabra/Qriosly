import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';
const errorCount = new Counter('custom_errors');
const BASE_URL = 'http://localhost:8080';
const QUIZ_ID = __ENV.QUIZ_ID;
const QUESTION_ID = __ENV.QUESTION_ID;
const OPTION_ID = __ENV.OPTION_ID;
const TARGET_VUS = parseInt(__ENV.TARGET_VUS || 50);
export const options = {
      scenarios: {
            quiz_flow: {
                  executor: 'ramping-vus',
                  startVUs: 0,
                  stages: [
                        { duration: '30s', target: TARGET_VUS },
                        { duration: '1m', target: TARGET_VUS },
                        { duration: '30s', target: 0 },
                  ],
            },
      },
      thresholds: {
            http_req_failed: ['rate<0.01'],
            http_req_duration: ['p(95)<500'],
      },
};
let authHeaders = null;
function ensureLoggedIn() {
      if (authHeaders) return;
      const uniqueId = `${__VU}_${Date.now()}`;
      const email = `student_${uniqueId}@loadtest.com`;
      const password = 'loadtest123';

      const registerRes = http.post(`${BASE_URL}/register`, JSON.stringify({
            name: `Load Test ${uniqueId}`,
            email: email,
            password: password,
            role: 'student',
      }), { headers: { 'Content-Type': 'application/json' }, tags: { name: 'Register' } });
      check(registerRes, { 'register succeeded': (r) => r.status === 201 }) || errorCount.add(1);
      const loginRes = http.post(`${BASE_URL}/login`, JSON.stringify({
            email: email,
            password: password,
      }), { headers: { 'Content-Type': 'application/json' }, tags: { name: 'Login' } });
      check(loginRes, { 'login succeeded': (r) => r.status === 200 });
      const token = loginRes.json('token');
      authHeaders = { headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' } };
}

export default function () {
      ensureLoggedIn();
      sleep(0.5);
      const quizRes = http.get(`${BASE_URL}/quizzes/${QUIZ_ID}`, { ...authHeaders, tags: { name: 'GetQuiz' } });
      check(quizRes, { 'fetch quiz succeeded': (r) => r.status === 200 }) || errorCount.add(1);
      const questionsRes = http.get(`${BASE_URL}/quizzes/${QUIZ_ID}/questions`, { ...authHeaders, tags: { name: 'GetQuestions' } });
      check(questionsRes, { 'fetch questions succeeded': (r) => r.status === 200 }) || errorCount.add(1);
      const optionRes = http.get(`${BASE_URL}/questions/${QUESTION_ID}/options`, { ...authHeaders, tags: { name: 'GetOptions' } });
      check(optionRes, { 'fetch options succeeded': (r) => r.status === 200 }) || errorCount.add(1);
      sleep(0.5);
      const attemptRes = http.post(`${BASE_URL}/quizzes/${QUIZ_ID}/attempts`, null, { ...authHeaders, tags: { name: 'StartAttempt' } });
      const attemptOk = check(attemptRes, { 'start attempt success': (r) => r.status === 201 });
      if (!attemptOk) { errorCount.add(1); return; }
      const attemptId = attemptRes.json('id');
      sleep(1);
      const answerRes = http.post(`${BASE_URL}/attempts/${attemptId}/answers`, JSON.stringify({
            question_id: QUESTION_ID,
            selected_option_id: OPTION_ID,
      }), { ...authHeaders, tags: { name: 'SubmitAnswer' } });
      check(answerRes, { 'submit answer succeeded': (r) => r.status === 201 }) || errorCount.add(1);
      sleep(0.5);
      const submitRes = http.post(`${BASE_URL}/attempts/${attemptId}/submit`, null, { ...authHeaders, tags: { name: 'SubmitAttempt' } });
      const submitOk = check(submitRes, {
            'submit attempt succeeded': (r) => r.status === 200,
            'score is present': (r) => r.json('score') !== null,
      });
      if (!submitOk) { errorCount.add(1); }
}