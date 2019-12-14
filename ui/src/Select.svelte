<script>
	import { onMount } from 'svelte';
    export let selected;
    let exportElement;

    onMount(() => {
        exportElement.selectedIndex = selected;
    });
</script>

<select id="input" class="select" bind:value={selected} bind:this={exportElement}>
		<option value=0>Health</option>
		<option value=1>Team</option>
</select>

<style>
    select {
        display:inline-block;
        width: 100px;
    }
</style>