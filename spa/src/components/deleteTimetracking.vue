<template>
  <div class="q-pa-md">
    <q-icon @click="dialog = true" class="float-right" color="red" name="delete" size="24px"/>

    <q-dialog v-model="dialog">
      <q-card>
        <q-card-section class="row items-center">
          <h6 class="no-margin">Are you sure you want to delete <br/>
            your tracked time with ID {{this.Time.id}}</h6>
        </q-card-section>

        <!-- Notice v-close-popup -->
        <q-card-actions align="right">
          <q-btn color="primary" flat label="Cancel" v-close-popup="cancelEnabled"/>
          <q-btn @click="deleteTime()" color="primary" flat label="Delete" v-close-popup/>
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script>
export default {
  props: ['Time'],
  data () {
    return {
      dialog: false,
      cancelEnabled: true
    }
  },
  methods: {
    deleteTime () {
      // console.log('time: ' + JSON.stringify(this.Time))
      this.$store.dispatch('vault/deleteTime', this.Time.id)
    }
  }
}
</script>
