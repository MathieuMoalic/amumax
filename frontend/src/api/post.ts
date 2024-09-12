import { setAlert } from "$lib/alerts/alert";

export async function post(endpoint: string, data: any) {
    const response = await fetch(`/api/${endpoint}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ ...data })
    });

    if (!response.ok) {
        try {
            const errorData = await response.json();
            setAlert(errorData.error);
        } catch (e) {
            setAlert("An error occurred");
        }
    }
}
