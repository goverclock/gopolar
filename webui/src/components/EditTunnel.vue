<template>
    <el-dialog title="Edit Tunnel" @open="onOpen">
        <el-form ref="formRef" :model="form" :rules="rules">
            <el-form-item label="Name" prop="name" label-width="140px">
                <el-input v-model="form.name" maxlength="21" show-word-limit placeholder="New tunnel" type="text" clearable
                    style="margin-right: 100px;"></el-input>
            </el-form-item>
            <el-form-item label="Souce" prop="source_port" label-width="140px">
                <el-text style="margin-right: 10px;">localhost:</el-text>
                <el-input-number v-model="form.source_port" :min="0" :max="65535" :value-on-clear="0" />
            </el-form-item>
            <el-form-item label="Dest" prop="dest" label-width="140px">
                <el-input v-model="form.dest" maxlength="21" placeholder="e.g. 1.1.1.1:1234" type="text" clearable
                    style="margin-right: 100px;"></el-input>
            </el-form-item>
        </el-form>

        <template #footer>
            <span class="dialog-footer">
                <el-button type="primary" @click="handleConfirm">Confirm</el-button>
                <el-button @click="handleCancel">Cancel</el-button>
            </span>
        </template>
    </el-dialog>
</template>

<script setup>
import { EditTunnelReq } from '../request/api'

const props = defineProps(['tunnel'])
const emit = defineEmits(['off', 'refresh'])

const formRef = ref();
const form = reactive({
    name: '',
    source_port: 1024,
    dest: '',
})
const rules = computed(() => {
    return {
        source_port: {
            required: true,
            message: "source port is required",
            trigger: ["change", "blur"],
        },
        dest: {
            required: true,
            trigger: ["change", "blur"],
            validator: (rule, value, callback) => {
                if (/^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?):(6553[0-5]|655[0-2][0-9]|65[0-4][0-9][0-9]|6[0-4][0-9][0-9][0-9][0-9]|[1-5](\d){4}|[1-9](\d){0,3})$/.test(value)) {
                    callback()
                    return;
                }
                // may contain "localhost"
                let sp = value.split(":")
                if (sp.length != 2) {
                    callback(new Error("invalid address"))
                    return;
                }
                if (sp[0] != "localhost") {
                    callback(new Error("invalid address"))
                    return;
                }
                let port = Number(sp[1])
                if (port <= 0 || port > 65535) {
                    callback(new Error("invalid address"))
                    return;
                }
                callback()
            }
        },
    }
});

function onOpen() {
    form.name = props.tunnel.name
    form.source_port = Number(props.tunnel.source.split(":")[1])
    form.dest = props.tunnel.dest
}

function handleConfirm() {
    formRef.value.validate(valid => {
        if (!valid) return

        let name = form.name
        if (name == '') {
            name = "New tunnel"
        }
        EditTunnelReq(props.tunnel.id, name, "localhost:" + form.source_port, form.dest)
            .then(res => {
                ElMessage({ message: `Tunnel ${name}(ID=${props.tunnel.id}) updated`, type: "success" })
                emit("refresh") // TODO: update single tunnel
                emit("off")
            }).catch((e) => { console.error(e) })
    })
}

function handleCancel() {
    emit("off")
}

</script>
