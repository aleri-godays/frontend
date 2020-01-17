<template>
  <q-page class="q-pa-lg">
    <h5 class="q-mt-none">Your Projects:</h5>

    <q-list bordered class="rounded-borders">
      <q-expansion-item
        :key="l.ID" expand-separator
        v-for="l in myProjects"
      >
        <template v-slot:header>

          <q-item-section>
            {{l.ID}} {{l.Name}}
          </q-item-section>

          <q-item-section>
            Total Duration: {{countTotal(l.ID)}}
          </q-item-section>

          <q-item-section side>
            <Timetracker :Project="l"/>
          </q-item-section>
        </template>
        <q-card
          :key="e.ID"
          v-for="e in myEntries">
          <q-card-section
            v-if="e.project_id===l.ID">
            <TimeDelete class="float-right" :Time="e"/>
            <TimeEditor class="float-right" :Time="e"/>
            ID: {{e.id}} <br/>
            Date: {{dateformatter(e.date)}}<br/>
            Comment: {{e.comment}}<br/>
            Single duration: {{timeConvert(e.duration)}}
          </q-card-section>
        </q-card>
      </q-expansion-item>
    </q-list>
  </q-page>
</template>

<script>
import { mapGetters } from 'vuex'
import moment from 'moment'
import Timetracker from '../components/addTimetracking'
import TimeEditor from '../components/editTimetracking'
import TimeDelete from '../components/deleteTimetracking'

export default {
  name: 'PageIndex',
  components: {
    Timetracker,
    TimeEditor,
    TimeDelete
  },
  data () {
    return {
      leftDrawerOpen: false
    }
  },
  created () {
    // this.$store.dispatch('vault/loadSample')
    this.$store.dispatch('vault/getProjects')
    this.$store.dispatch('vault/getEntries')
  },
  computed: {
    ...mapGetters({
      myProjects: 'vault/getAllProjects',
      myEntries: 'vault/getAllEntries'
    })
  },
  methods: {
    dateformatter (value) {
      if (value) {
        return moment(String(value)).format('MM/DD/YYYY hh:mm')
      }
    },

    todaysDate (value) {
      if (value !== 0) {
        return moment(value).toDate.toString
      } else return 0
    },

    float2int (value) {
      return value | 0
    },
    timeConvert (n) {
      var se = n
      var num = (se / 60)
      var hours = (num / 60)
      var rhours = Math.floor(hours)
      var minutes = (hours - rhours) * 60
      var rminutes = Math.round(minutes)

      var days = (hours / 24)
      var rdays = Math.floor(days)
      var h = (days - rdays) * 24
      var rh = Math.round(h)
      return rdays + ' day(s) ' + rh + ':' + rminutes + ' h'
    },

    countTotal (id) {
      var duration = 0
      for (var i = 0; i < this.myEntries.length; i++) {
        if (this.myEntries[i].project_id === id) {
          duration = duration + this.myEntries[i].duration
        }
      }
      // if (duration !== 0) {
      //   var t = new Date(duration)
      //   var formatted = moment(t).format('HH:mm')
      //   return formatted
      // } else return 0
      // return moment.duration(duration).asHours()
      if (duration !== 0) {
        return this.timeConvert(duration)
      } else return 0
    }
  }
}
</script>
