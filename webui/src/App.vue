<template>
  <div class="header-cont">
    <img src="./assets/gopolar.jpg">
    <h1>
      gopolar
    </h1>
  </div>
  <el-text>version {{ about['version'] }}</el-text>
  <br>
  <br>

  <el-button link type="info" :icon="Plus" @click="createViewVisible = true">new tunnel</el-button>
  <el-button link type="info" :icon="Refresh" @click="handleRefresh">refresh</el-button>

  <el-table :data="tableData">
    <el-table-column align="center" prop="id" label="ID" width="100px" />
    <el-table-column prop="name" label="Name" />
    <el-table-column prop="source" label="Source" />
    <el-table-column prop="dest" label="Dest" />
    <el-table-column align="center" prop="status" label="Status" width="100px">
      <template v-slot="{ row }">
        <el-tag type="success" v-if="row.status">RUNNING</el-tag>
        <el-tag type="info" v-else>STOPPED</el-tag>
      </template>
    </el-table-column>
    <el-table-column align="center" fixed="right" label="Operations" width="250px">
      <template v-slot="{ row }">
        <el-button link type="primary" size="small" @click="handleToggle(row)">Toggle</el-button>
        <el-button link type="primary" size="small" @click="selectedTunnel = row; editViewVisible = true">Edit</el-button>
        <el-button link type="danger" size="small" @click="handleDelete(row)">Delete</el-button>
      </template>
    </el-table-column>
  </el-table>

  <CreateTunnel v-model="createViewVisible" @off="createViewVisible = false" @refresh="handleRefresh" />
  <EditTunnel v-model="editViewVisible" :tunnel="selectedTunnel" @off="editViewVisible = false"
    @refresh="handleRefresh" />
</template>

<script setup>
import { Refresh, Plus } from '@element-plus/icons-vue'
import { ElTag } from 'element-plus';
import { AboutReq, DeleteTunnelReq, ToggleTunnelReq, GetTunnelListReq } from './request/api'
import CreateTunnel from './components/CreateTunnel.vue'
import EditTunnel from './components/EditTunnel.vue'

const about = ref({})
AboutReq().then(res => {
  about.value = res.data['about']
}).catch((e) => { console.error(e) })

// init tunnel table
const tableData = ref([])
function handleRefresh() {
  tableData.value = []
  GetTunnelListReq().then(res => {
    let tunnels = res.data["tunnels"]
    for (let t of tunnels) {
      tableData.value.push({
        id: t.id,
        name: t.name,
        source: t.source,
        dest: t.dest,
        status: t.enable,
      })
    }
  }).catch(e => { console.error(e) })
}
handleRefresh()

function handleToggle(row) {
  row.status = !row.status
  ToggleTunnelReq(row.id).then(res => {
    if (row.status) {
      ElMessage({ message: `Tunnel ${row.name}(ID=${row.id}) started`, type: "success" })
    } else {
      ElMessage({ message: `Tunnel ${row.name}(ID=${row.id}) stopped`, type: "success" })
    }

    // TODO(pending): need new API for gopolar core: GetTunnelInfo(id)
    // or else the whole table blinks
    handleRefresh()
  }).catch(e => { console.error(e) })
}

function handleDelete(row) {
  DeleteTunnelReq(row.id).then(res => {
    ElMessage({ message: `Tunnel ${row.name}(ID=${row.id}) deleted`, type: "success" })
    handleRefresh()
  }).catch((e) => { console.error(e) })
}

const createViewVisible = ref(false)
const editViewVisible = ref(false)
const selectedTunnel = ref({})

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