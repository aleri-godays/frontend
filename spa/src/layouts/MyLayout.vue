<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn v-if="User.loginState"
          @click="leftDrawerOpen = !leftDrawerOpen"
          aria-label="Menu"
          dense
          flat
          icon="menu"
          round
        />

        <q-toolbar-title>
          Welcome to the Godays 2020
        </q-toolbar-title>
        <q-btn v-if="this.$store.state.vault.user.loginState"
               @click="logout"
               class="absolute-right"
               flat
               label="logout"
        />
        <q-btn v-else
          @click="login"
          class="absolute-right"
          flat
          label="login"
        />
      </q-toolbar>
    </q-header>

    <q-drawer v-if="User.loginState"
    v-model="leftDrawerOpen"
    show-if-above
    :width="200"
    :breakpoint="400"
    >
    <q-scroll-area style="height: calc(100% - 150px); margin-top: 50px; border-right: 1px solid #ddd">
      <q-list padding>

        <q-item
          to="/projects"
          clickable
          v-ripple>
          <q-item-section avatar>
            <q-icon name="work" />
          </q-item-section>

          <q-item-section>
            Projects
          </q-item-section>
        </q-item>
      </q-list>
    </q-scroll-area>

    <q-img class="absolute-top" src="https://cdn.quasar.dev/img/material.png" style="height: 50px">
      <div class="bg-transparent">
        <div class="text-weight-bold">{{User.name}}</div>
      </div>
    </q-img>
    </q-drawer>

    <q-page-container>
      <keep-alive>
      <router-view/>
      </keep-alive>
    </q-page-container>
  </q-layout>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'MyLayout',
  data () {
    return {
      leftDrawerOpen: false
    }
  },
  methods: {
    login () {
      this.$store.dispatch('vault/login')
    },
    logout () {
      this.$store.dispatch('vault/logout')
    }
  },
  computed: {
    ...mapGetters({
      User: 'vault/getUser'
    })
  }

}
</script>
