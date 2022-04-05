<script>
    import Line from "svelte-chartjs/src/Line.svelte";

    import { onMount } from "svelte";
    import { getTables } from "../api/crud";

    let xlabel = "";
    let ylabel = "";
    let tables = [];

    let chart = {
        labels: [],
        datasets: [
            {
                label: "My First dataset",
                fill: false,
                lineTension: 0.1,
                backgroundColor: "rgba(225, 204,230, .3)",
                borderColor: "rgb(205, 130, 158)",
                // borderCapStyle: "butt",
                // borderDash: [],
                // borderDashOffset: 0.0,
                // borderJoinStyle: "miter",
                pointBorderColor: "rgb(205, 130,1 58)",
                pointBackgroundColor: "rgb(255, 255, 255)",
                pointBorderWidth: 10,
                pointHoverRadius: 5,
                pointHoverBackgroundColor: "rgb(0, 0, 0)",
                pointHoverBorderColor: "rgba(220, 220, 220,1)",
                pointHoverBorderWidth: 2,
                pointRadius: 1,
                pointHitRadius: 10,
                pointStyle: "cross",
                data: [],
            },
        ],
    };
    $: chart.labels = tables[xlabel];
    $: chart.datasets[0].data = tables[ylabel];

    onMount(async () => {
        const res = await getTables();
        tables = res.data;
        xlabel = "mz";
        ylabel = "mz";
    });
</script>

<main>
    <div class="ybtns">
        {#each Object.entries(tables) as [name, _]}
            <label>
                <input
                    type="radio"
                    id={name}
                    bind:group={ylabel}
                    name="ylabel"
                    value={name}
                />
                <span>{name}</span>
            </label>
        {/each}
    </div>
    <div class="c1">
        <div class="plot">
            <Line data={chart} options={{ responsive: true }} />
        </div>
        <div class="xbtns">
            {#each Object.entries(tables) as [name, _]}
                <label>
                    <input
                        type="radio"
                        id={name}
                        bind:group={xlabel}
                        name="xlabel"
                        value={name}
                    />
                    <span>{name}</span>
                </label>
            {/each}
        </div>
    </div>
</main>

<style>
    main {
        width: 100%;
        height: 100%;
        border: 2px solid blue;
        display: flex;
    }
    .ybtns {
        width: 10%;
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
    }
    .c1 {
        width: 90%;
        height: 100%;
        display: flex;
        flex-direction: column;
    }
    .plot {
        width: 100%;
        height: 90%;
    }
    .xbtns {
        width: 100%;
        height: 10%;
        display: flex;
    }
    label {
        display: flex;
        cursor: pointer;
        font-weight: 500;
        position: relative;
        overflow: hidden;
    }
    input {
        position: absolute;
        left: -9999px;
    }
    input:checked + span {
        background-color: var(--btn);
    }

    span {
        display: flex;
        vertical-align: middle;
        align-items: center;
        width: 4rem;
        height: 2rem;
        margin: 0.3rem;
        padding: 0.375em 0.75em 0.375em 0.375em;
        border-style: solid;
        border-radius: 0.2rem;
        border-width: 0.1rem;
        border-color: var(--btn);
        transition: 0.25s ease;
        font-size: 0.6rem;
    }
    span:hover {
        background-color: var(--hover);
    }
</style>
