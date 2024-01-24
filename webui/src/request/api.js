import request from "./request";

export function GetTunnelListReq() {
    return request({
        url: 'tunnels/list',
        method: 'get',
    })
}

export function CreateTunnelReq(name, source, dest) {
    return request({
        url: 'tunnels/create',
        method: 'post',
        data: {
            name,
            source,
            dest,
        }
    })
}

export function EditTunnelReq(id, name, source, dest) {
    return request({
        url: 'tunnels/edit/' + id,
        method: 'post',
        data: {
            name,
            source,
            dest,
        }
    })
}

export function ToggleTunnelReq(id) {
    return request({
        url: 'tunnels/toggle/' + id,
        method: 'post',
    })
}

export function DeleteTunnelReq(id) {
    return request({
        url: 'tunnels/delete/' + id,
        method: 'delete',
    })
}

export function AboutReq() {
    return request({
        url: 'about',
        method: 'get',
    })
}
