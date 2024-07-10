import { post } from "$api/post";


export function postCommand(command: string) {
    post('console-command', { command });
}
