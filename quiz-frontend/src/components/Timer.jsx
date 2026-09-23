export default function Timer({ remainingMs }) {
      if (remainingMs === null) return null;

      const expired = remainingMs <= 0;
      const totalSeconds = Math.max(0, Math.floor(remainingMs / 1000));
      const minutes = Math.floor(totalSeconds / 60);
      const seconds = totalSeconds % 60;
      const formatted = `${minutes}:${seconds.toString().padStart(2, '0')}`;

      return (
            <span className={`font-mono text-sm font-semibold ${expired ? 'text-red-500' : 'text-slate-700'}`}>
                  ⏱ {expired ? 'Time Expired' : formatted}
            </span>
      );
}