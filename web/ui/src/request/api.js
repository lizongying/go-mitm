import axios from "axios";

const host = 'http://localhost:8083'
const action = params => {
    return axios.get(host + '/action', {
        params,
    });
};

const info = () => {
    return axios.get(host + '/info')
};

const event = () => {
    let url = host + '/event'
    return new EventSource(url);
};

export {action, info, event}