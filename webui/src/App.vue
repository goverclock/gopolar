<template>
  <!-- TODO: gopolar logo here -->
  <TunnelList></TunnelList>
  <!-- <HelloWorld msg="Vite + Vue" /> -->
</template>

<script setup>
import TunnelList from './components/TunnelList.vue'
import { AboutReq, DeleteTunnelReq, ToggleTunnelReq, EditTunnelReq, CreateTunnelReq, GetTunnelListReq } from './request/api'

GetTunnelListReq().then(res => {
  console.log(res)
}).catch((e) => { console.error(e) })

CreateTunnelReq("fromwebui", "localhost:2332", "localhost:123")
  .then(res => {
    console.log(res)
  }).catch((e) => { console.error(e) })

EditTunnelReq(2, "hehe", "localhost:231", "1.1.1.1:3")
  .then(res => {
    console.log(res)
  }).catch((e) => { console.error(e) })

ToggleTunnelReq(3).then(res => {
  console.log(res)
}).catch((e) => { console.error(e) })

DeleteTunnelReq(3).then(res => {
  console.log(res)
}).catch((e) => { console.error(e) })

AboutReq().then(res => {
  console.log(res)
}).catch((e) => { console.error(e) })

</script>
