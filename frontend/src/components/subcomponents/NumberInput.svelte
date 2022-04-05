<script>
    export let value = 0;
    export let label = "label";
</script>

<main>
    <label for="input">{label}</label>
    <input type="number" id="input" name="quantity" step="any" bind:value />
</main>

<style>
    input {
        background-color: blue;
    }
    /* removes arrows */
    input::-webkit-inner-spin-button {
        -webkit-appearance: none;
        margin: 0;
    }
    input[type="number"] {
        -moz-appearance: textfield;
    }
</style>
