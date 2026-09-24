export default function QuestionPalette({ questions, answers, currentIndex, onSelect }) {
      return (
            <div className="flex flex-wrap gap-2">
                  {questions.map((q, idx) => {
                        const isAnswered = !!answers[q.id];
                        const isCurrent = idx === currentIndex;

                        return (
                              <button
                                    key={q.id}
                                    onClick={() => onSelect(idx)}
                                    className={`w-10 h-10 rounded-lg text-sm font-medium border flex items-center justify-center transition-colors ${isCurrent
                                                ? 'bg-indigo-700 text-white border-indigo-700'
                                                : isAnswered
                                                      ? 'bg-emerald-50 text-emerald-600 border-emerald-300'
                                                      : 'bg-white text-slate-500 border-slate-200'
                                          }`}
                              >
                                    {idx + 1}
                              </button>
                        );
                  })}
            </div>
      );
}