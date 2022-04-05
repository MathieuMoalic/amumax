import Api from "./Api"


export const getDir = async () => {
    try {
        const response = await Api.get("/dir");
        return response;
    } catch (error) {
        console.error(error);
    }
};
export const getTables = async () => {
    try {
        const response = await Api.get("/tables");
        return response;
    } catch (error) {
        console.error(error);
    }
};
export const getImage = async () => {
    try {
        const response = await Api.get("/image");
        return response;
    } catch (error) {
        console.error(error);
    }
};

// export const getWord = async (word_id) => {
//     try {
//         const response = await Api.get(`/word/${word_id}/`);
//         return response;
//     } catch (error) {
//         console.error(error);
//     }
// };

// export const getRandomWord = async () => {
//     try {
//         const response = await Api.get(`/random_word/`);
//         return response;
//     } catch (error) {
//         console.error(error);
//     }
// };

// export const postWord = async (word) => {
//     try {
//         const response = await Api.post(`/word/`, word);
//         return response;
//     } catch (error) {
//         console.error(error);
//     }
// };

// export const updateWord = async (word_id, word) => {
//     try {
//         const response = await Api.put(`/word/${word_id}`, word);
//         return response;
//     } catch (error) {
//         console.error(error);
//     }
// };