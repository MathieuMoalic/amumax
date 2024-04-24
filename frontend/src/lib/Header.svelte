<script lang="ts">
	let title = 'mytest/t2';
	let progress = 0;
	let status = 'Paused';
</script>

<span>
	<span>{title}</span> &nbsp;
</span>
<progress max="100" value="0"></progress>
<span id="busy">Paused</span>
<span class="ErrorBox"></span>
<span id="MsgBox" class="ErrorBox"></span>
<br />

<style>
	.ErrorBox {
		color: red;
		font-weight: bold;
		font-size: 1.5em;
	}
	span {
		color: rgb(255, 255, 255);
		font-weight: bold;
		font-size: 1.5em;
	}
</style>
