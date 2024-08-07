export async function post(endpoint: string, data: any) {
    const response = await fetch(`/api/${endpoint}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ ...data })
    });
    if (!response.ok) {
        throw new Error(`Failed to post ${endpoint}`);
    }
}