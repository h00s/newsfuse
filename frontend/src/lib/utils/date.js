export function formatDate(date) {
  date = new Date(date);
  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');

  return `${day}.${month}.${year}. ${hours}:${minutes}`;
}

export function humanizeDuration(date) {
  date = new Date(date);
  const now = new Date();
  const diff = now - date;

  const seconds = Math.floor(diff / 1000);
  if (seconds < 60) {
    return `prije ${seconds} sek`;
  }

  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) {
    return `prije ${minutes} min`;
  }

  const hours = Math.floor(minutes / 60);
  if (hours < 24) {
    return `prije ${hours} h`;
  }

  const days = Math.floor(hours / 24);
  return `prije ${days} dan(a)`;
}