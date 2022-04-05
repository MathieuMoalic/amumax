<script>
    import NumberInput from "./subcomponents/NumberInput.svelte";
    import { Meshx } from "../stores.js";
</script>

<main>
    <h5>Mesh</h5>
    <NumberInput label="Meshx" value={$Meshx} />
</main>
