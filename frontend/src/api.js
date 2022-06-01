let api = "https://groceriesapi.matmoa.xyz";

let addItem = async () => {
    if (newWord != "") {
        let res = await fetch(api + "/" + newWord, { method: "post" });
        items = await res.json();
        newWord = "";
    }
};

let removeItem = async (id) => {
    let res = await fetch(api + "/" + id, { method: "delete" });
    items = await res.json();
};
