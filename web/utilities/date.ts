
export function formatUnixToLocalDateTime(unixTimestamp: number): string {
  const date = new Date(unixTimestamp * 1000); // Convert from seconds to milliseconds

  const day = date.getDate();
  const monthNames = [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December'
  ];
  const month = monthNames[date.getMonth()];
  const year = date.getFullYear();
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');

  return `${day} ${month} ${year} ${hours}:${minutes}`;
}


export function formatMilliseconds(ms: number): string {
  if (ms >= 3600000) {
    return '--:--:---';
  }

  const minutes = Math.floor(ms / (1000 * 60));
  const seconds = Math.floor((ms % (1000 * 60)) / 1000);
  const milliseconds = ms % 1000;

  const pad = (num: number, size: number) => num.toString().padStart(size, '0');

  return `${pad(minutes, 2)}:${pad(seconds, 2)}:${pad(milliseconds, 3)}`;
}
