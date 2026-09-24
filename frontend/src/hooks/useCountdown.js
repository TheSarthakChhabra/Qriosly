import { useState, useEffect } from 'react';

export function useCountdown(deadline) {
      const [remainingMs, setRemainingMs] = useState(() => (deadline ? deadline - Date.now() : null));

      useEffect(() => {
            if (!deadline) return;

            const interval = setInterval(() => {
                  setRemainingMs(deadline - Date.now());
            }, 1000);

            return () => clearInterval(interval);
      }, [deadline]);

      return remainingMs;
}