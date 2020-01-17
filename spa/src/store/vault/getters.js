export function oddLines (state) {
  return state.lines.all.filter(item => (item.ID % 2 === 0))
}

export function getAllProjects (state) {
  return state.projects.all
}

export function getAllEntries (state) {
  return state.entries.all
}

export function getUser (state) {
  return state.user
}
