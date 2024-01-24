import axios from "axios";
import { ElMessage } from "element-plus";

export default function request(options) {
    return new Promise((resolve, reject) => {
        const instance = axios.create({
            baseURL: "http://localhost:8080",
            headers: {
                'Content-Type': 'application/json',
            },
            timeout: 100000,
            responseType: 'json',
            ...options
        });
        // response 响应拦截器
        instance.interceptors.response.use(
            (response) => {
                return response.data;
            },
            (err) => {
                if (err && err.response) {
                    switch (err.response.status) {
                        case 400:
                            err.message = "请求错误";
                            break;
                        case 401:
                            err.message = "未授权，请登录";
                            break;
                        case 403:
                            err.message = "拒绝访问";
                            break;
                        case 404:
                            err.message = `请求地址出错: ${err.response.config.url}`;
                            break;
                        case 408:
                            err.message = "请求超时";
                            break;
                        case 500:
                            err.message = "服务器内部错误";
                            break;
                        case 501:
                            err.message = "服务未实现";
                            break;
                        case 502:
                            err.message = "网关错误";
                            break;
                        case 503:
                            err.message = "服务不可用";
                            break;
                        case 504:
                            err.message = "网关超时";
                            break;
                        case 505:
                            err.message = "HTTP版本不受支持";
                            break;
                        default:
                    }
                }

                if (err.message) {
                    ElMessage({ message: err.message, type: "error" });
                }
                return Promise.reject(err); // 返回接口返回的错误信息
            }
        );

        // 请求处理
        instance(options)
            .then((res) => {
                if (res.success) {
                    resolve(res);
                } else {
                    ElMessage({ message: res.err_msg || "操作失败", type: "error", showClose: false });
                    reject(res);
                }
            })
            .catch((error) => {
                reject(error);
            });
    });
}
