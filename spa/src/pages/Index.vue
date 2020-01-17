<template>
  <q-page class="q-pa-lg">
    <h5 class="q-mt-none">Startpage</h5>
  </q-page>
</template>

<script>
// if user is not logged in -> say "you are not logged in ! Please do so "
import { mapState, mapGetters } from 'vuex'
export default {
  name: 'PageIndex',
  data () {
    return {
      leftDrawerOpen: false
    }
  },
  created () {
    this.$store.dispatch('vault/loadSample')
    this.$store.dispatch('vault/authorize')
    // this.$store.dispatch('vault/getProjects')
  },
  computed: {
    ...mapState({
      myLines: state => state.vault.lines.all
    }),
    ...mapGetters({
      myOddLines: 'vault/oddLines'
    }),
    ...mapGetters({
      myProjects: 'vault/getAllProjects'
    })
  }
}
</script>
