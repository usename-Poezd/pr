const KEY='pulse-created-polls'
export function getCreatedPolls(){try{return JSON.parse(localStorage.getItem(KEY)||'[]')}catch{return []}}
export function saveCreatedPoll(poll){const polls=getCreatedPolls().filter(item=>item.id!==poll.id);localStorage.setItem(KEY,JSON.stringify([poll,...polls]))}
