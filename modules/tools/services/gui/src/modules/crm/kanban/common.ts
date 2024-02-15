function createDate(date: Date) {
    return date.setMinutes(date.getMinutes() + date.getTimezoneOffset());
  }
  

export function formatDateWithTime(date: string) {
    return date === '0001-01-01T00:00:00Z'
      ? ''
      : new Intl.DateTimeFormat('ua', {
          year: 'numeric',
          month: 'numeric',
          day: 'numeric',
          hour: 'numeric',
          minute: 'numeric',
          hour12: false,
        })
          .format(createDate(new Date(date)))
          .replace(',', '');
  }
  
  export function formatDate(date: string) {
    return date === '0001-01-01T00:00:00Z'
      ? ''
      : new Intl.DateTimeFormat('ua', {
          year: 'numeric',
          month: 'numeric',
          day: 'numeric',
          hour12: false,
        }).format(createDate(new Date(date)));
  }
  
  export function formatTime(date: string) {
    return new Intl.DateTimeFormat('ua', {
      hour: 'numeric',
      minute: 'numeric',
      hour12: false,
    }).format(createDate(new Date(date)));
  }
  
  export function formatTimeWithHoursMinutes(date: string) {
    return new Intl.DateTimeFormat('ua', {
      hour: 'numeric',
      minute: 'numeric',
      hour12: false,
    }).format(createDate(new Date(date)));
  }
  