<template>
  <el-row class="header-cont">
    <img src="./assets/gopolar.jpg">
    <h1>
      gopolar
    </h1>
  </el-row>

  <!-- <TunnelList></TunnelList> -->
  <!-- <HelloWorld msg="Vite + Vue" /> -->
  <el-table :data="tableData" style="width: 100%">
    <el-table-column prop="date" label="Date" />
    <el-table-column prop="name" label="Name" />
    <el-table-column prop="address" label="Address" />
  </el-table>
</template>

<script setup>
import TunnelList from './components/TunnelList.vue'
import { AboutReq, DeleteTunnelReq, ToggleTunnelReq, EditTunnelReq, CreateTunnelReq, GetTunnelListReq } from './request/api'

const tableData = [
  {
    date: '2016-05-03',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-02',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-04',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-01',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
]

// GetTunnelListReq().then(res => {
//   console.log(res)
// }).catch((e) => { console.error(e) })

// CreateTunnelReq("fromwebui", "localhost:2332", "localhost:123")
//   .then(res => {
//     console.log(res)
//   }).catch((e) => { console.error(e) })

// EditTunnelReq(2, "hehe", "localhost:231", "1.1.1.1:3")
//   .then(res => {
//     console.log(res)
//   }).catch((e) => { console.error(e) })

// ToggleTunnelReq(3).then(res => {
//   console.log(res)
// }).catch((e) => { console.error(e) })

// DeleteTunnelReq(3).then(res => {
//   console.log(res)
// }).catch((e) => { console.error(e) })

AboutReq().then(res => {
  console.log(res)
}).catch((e) => { console.error(e) })

</script>

<style>
.header-cont {
  /* background-color: #f5f7f9; */
  display: flex;
  justify-content: center;
  height: 100%;
  padding-top: 0;
  padding-bottom: 10px;

  h1 {
    margin: 0;
    font-size: 40px;
  }

  img {
    max-height: 50px;
    max-width: 50px;
    margin: 0 0 0 0;
    display: flex;
    justify-content: space-between;
  }

  .gap {
    margin-right: 20px;
  }
}
</style>