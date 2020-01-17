export function setLines (state, items) {
  if (items === null) {
    items = []
  }
  state.lines.all = items
}
export function setUser (state, usr) {
  if (usr === null) {
    usr = ''
  }
  state.user.name = usr.login
  state.user.token = usr.jwt
  state.user.loginState = true
}

export function setProjects (state, items) {
  if (items === null) {
    items = []
  }
  // console.log('project items: ', items)
  state.projects.all = items
  // console.log('project state: ', state.projects.all)
}

export function setEntries (state, items) {
  if (items === null) {
    items = []
  }
  state.entries.all = items
}
