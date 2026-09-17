import { useEffect, useState } from 'react';

function App() {
  const [result, setResult] = useState('loading...');

  useEffect(() => {
    fetch('http://localhost:8080/health')
      .then((res) => res.json())
      .then((data) => setResult(JSON.stringify(data)))
      .catch((err) => setResult('ERROR: ' + err.message));
  }, []);

  return <div>{result}</div>;
}

export default App;