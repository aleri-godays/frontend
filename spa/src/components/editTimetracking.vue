<template>
  <div class="q-pa-md">
    <q-icon @click="dialog = true" color="blue" name="edit" size="24px"/>

    <q-dialog v-model="dialog">
      <q-card>
        <q-card-section class="row items-center">
          <h6 class="no-margin">Edit time schedule with ID {{this.Time.id}}</h6>
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
          <q-input outlined hint="Duration" type="time" v-model="timeToSubmit.duration"/>
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
import { mapGetters } from 'vuex'
import { date } from 'quasar'
export default {
  props: ['Time'],
  data () {
    return {
      options: [],
      dialog: false,
      cancelEnabled: true,
      timeToSubmit: {
        id: this.Time.id,
        date: this.dateformatter(this.Time.date),
        project_id: this.Time.project_id,
        comment: this.Time.comment,
        duration: this.timeformatter(this.Time.duration)
      }
    }
  },
  computed: {
    ...mapGetters({
      myProjects: 'vault/getAllProjects'
    })
  },

  methods: {
    closeDialog () {
      this.$refs.qDateProxy.hide()
      this.$refs.qTimeProxy.hide()
    },
    todaysDate () {
      let timestamp = Date.now()
      return date.formatDate(timestamp, 'DD/MM/YYYY HH:mm')
    },
    save () {
      if (this.timeToSubmit.date === this.dateformatter(this.Time.date)) {
        this.timeToSubmit.date = this.Time.date
      } else {
        this.timeToSubmit.date = date.formatDate(this.timeToSubmit.date, 'YYYY-MM-DDTHH:00:mmZ')
      }
      var h = String(this.timeToSubmit.duration).substr(0, 2)
      var m = String(this.timeToSubmit.duration).substr(-2)
      this.timeToSubmit.duration = (h * 60 * 60 + m * 60)
      // console.log('time to submit: ' + JSON.stringify(this.timeToSubmit))
      this.$store.dispatch('vault/editTime', this.timeToSubmit)
    },
    dateformatter (value) {
      if (value) {
        return date.formatDate(value, 'DD/MM/YYYY HH:mm')
      }
    },
    timeformatter (value) {
      return date.formatDate(value, 'HH:mm')
    }
  }
}
</script>
