<script>
    import { Meshx } from "../stores.js";
    // console.log(Meshx);
    $: v = $Meshx * 1;
</script>

<main>
    <h5>Parameters</h5>
    <h6>{v}</h6>
</main>
