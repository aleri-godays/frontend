import axios from 'axios'

export function loadSample ({ commit }) {
  const sd = [
    { ID: 1, Title: 'Line1', Value: 1 },
    { ID: 2, Title: 'Line2', Value: 2 },
    { ID: 3, Title: 'Line3', Value: 3 }
  ]
  commit('setLines', sd)
}

export function login ({ commit }) {
  // console.log('calling /me')
  axios
    .get('/api/v1/frontend/me')
    .then(
      response => {
        commit('setUser', response.data)
        // console.log('response: ', response)
      })
    .catch(error => {
      console.log('error, redirecting to /login', error)
      window.location = '/login'
    }
    )
}

export function authorize ({ commit }) {
  // console.log('calling /me')
  axios
    .get('/api/v1/frontend/me')
    .then(
      response => {
        commit('setUser', response.data)
        // console.log('response: ', response)
      })
    .catch(error => {
      console.log('error ', error)
    }
    )
}

export function getProjects ({ commit }) {
  axios
    .get('/api/v1/frontend/project')
    .then(
      response => {
        commit('setProjects', response.data)
        // console.log('project response: ', response)
      })
    .catch(
      (error) => console.log(error)
    )
}

export function getEntries ({ commit }) {
  axios
    .get('/api/v1/frontend/entry')
    .then(
      response => {
        commit('setEntries', response.data)
        // console.log('entries response: ', response)
      })
    .catch(
      (error) => console.log(error)
    )
}

export function logout ({ commit }) {
  axios
    .get('/logout')
    .then(
      response => {
        localStorage.clear()
        // console.log('response: ', response)
      })
    .then(
      window.location = '/'
    )
    .catch(
      (error) => console.log('error', error),
      window.location = '/logout',
      localStorage.clear()
    )
}

export function deleteTime ({ commit, dispatch }, timeIDToDelete) {
  axios
    .delete('/api/v1/frontend/entry/' + timeIDToDelete)
    .then(
      response => {
        // console.log('response: ', response)
        if (response) {
          dispatch('getEntries')
        }
      })
    .catch(error => {
      console.log('error, could not log out', error)
    }
    )
}

export function addTime ({ commit, dispatch }, entry) {
  // console.log('timeToSubmit: ' + JSON.stringify(entry))
  axios
    .post('/api/v1/frontend/entry', entry)
    .then(
      response => {
        // console.log('add time response: ', response)
        if (response) {
          dispatch('getEntries')
        }
      })
    .catch(error => {
      console.log('error, could not log out', error)
    }
    )
}

export function editTime ({ commit, dispatch }, entry) {
  // console.log('timeToSubmit: ' + JSON.stringify(entry))
  axios
    .put('/api/v1/frontend/entry/' + entry.id, entry)
    .then(
      response => {
        // console.log('response: ', response)
        if (response) {
          dispatch('getEntries')
        }
      })
    .catch(error => {
      console.log('error, could not log out', error)
    }
    )
}
