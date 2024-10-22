import { post } from "$api/post";


export function postResetError() {
    post('metrics/reset-error', { });
}
