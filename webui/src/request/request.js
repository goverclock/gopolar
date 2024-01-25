import axios from "axios";
import { ElMessage } from "element-plus";

export default function request(options) {
    return new Promise((resolve, reject) => {
        const instance = axios.create({
            baseURL: "http://" + window.location.hostname + ":7070",
            headers: {
                'Content-Type': 'application/json',
            },
            timeout: 100000,
            responseType: 'json',
            ...options
        });

        instance.interceptors.response.use(
            (response) => {
                return response.data;
            },
            (err) => {
                if (err.message) {
                    ElMessage({ message: err.message, type: "error" });
                }
                return Promise.reject(err);
            }
        );

        instance(options)
            .then((res) => {
                if (res.success) {
                    resolve(res);
                } else {
                    ElMessage({ message: res.err_msg || "Operation failed", type: "error", showClose: false });
                    reject(res);
                }
            })
            .catch((error) => {
                reject(error);
            });
    });
}
