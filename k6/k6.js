import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
    scenarios: {
        get: {
            executor: 'constant-vus',
            exec: 'get',
            vus: 500,
            duration: '500s',
        },
        post: {
            executor: 'constant-vus',
            exec: 'post',
            vus: 500,
            duration: '500s',
        },
    },
};


export function get() {
    // Make a GET request to the target URL
    http.get('http://localhost:3000/');

    // Sleep for 1 second to simulate real-world usage
    sleep(1);
}

export function post() {
    // Make a GET request to the target URL
    http.post('http://localhost:3000/post');

    // Sleep for 1 second to simulate real-world usage
    sleep(1.5);
}