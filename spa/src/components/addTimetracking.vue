<template>
  <div class="q-pa-md">
    <q-icon flat @click="dialog = true" color="green" name="add" size="24px"/>

    <q-dialog v-model="dialog">
      <q-card>
        <q-card-section class="row items-center">
          <h6 class="no-margin">Add a new time schedule to {{this.Project.Name}}</h6>
        </q-card-section>
        <q-card-section>
          <q-input outlined v-model="timeToSubmit.date" hint="Date">
            <template v-slot:prepend>
              <q-icon name="event" class="cursor-pointer">
                <q-popup-proxy transition-show="scale" transition-hide="scale" ref="qDateProxy">
                  <q-date v-model="timeToSubmit.date" mask="YYYY-MM-DD HH:mm" @input="closeDialog"/>
                </q-popup-proxy>
              </q-icon>
            </template>

            <template v-slot:append>
              <q-icon name="access_time" class="cursor-pointer">
                <q-popup-proxy transition-show="scale" transition-hide="scale" ref="qTimeProxy">
                  <q-time v-model="timeToSubmit.date" mask="YYYY-MM-DD HH:mm" format24h @input="closeDialog"/>
                </q-popup-proxy>
              </q-icon>
            </template>
          </q-input>
        </q-card-section>
        <q-card-section>
          <q-input outlined hint="Duration" label="hh:mm" mask="HH:mm" format24h type="time" v-model="timeToSubmit.duration"/>
        </q-card-section>
        <q-card-section>
          <q-input
            outlined
            hint="Add a comment"
            type="textarea"
            v-model="timeToSubmit.comment"
          />
        </q-card-section>

        <!-- Notice v-close-popup -->
        <q-card-actions align="right">
          <q-btn color="primary" flat label="Cancel" v-close-popup="cancelEnabled"/>
          <q-btn @click="save" color="primary" flat label="Save" v-close-popup/>
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script>
import { date } from 'quasar'
export default {
  props: ['Project'],
  data () {
    return {
      dialog: false,
      cancelEnabled: true,
      timeToSubmit: {
        date: '',
        project_id: this.Project.ID,
        comment: '',
        duration: ''
      }
    }
  },
  methods: {
    closeDialog () {
      this.$refs.qDateProxy.hide()
      this.$refs.qTimeProxy.hide()
    },
    save () {
      this.timeToSubmit.date = date.formatDate(this.timeToSubmit.date, 'YYYY-MM-DDTHH:00:mmZ')
      var h = String(this.timeToSubmit.duration).substr(0, 2)
      var m = String(this.timeToSubmit.duration).substr(-2)
      this.timeToSubmit.duration = (h * 60 * 60 + m * 60)
      // console.log('timeToSubmit: ' + JSON.stringify(this.timeToSubmit))
      this.$store.dispatch('vault/addTime', this.timeToSubmit)
    },
    todayDate () {
      let timeStamp = Date.now()
      this.timeToSubmit.date = date.formatDate(timeStamp, 'YYYY-MM-DD HH:mm')
    }
  },
  created () {
    this.todayDate()
  }
}
</script>
