<script lang="ts">
	import { headerState } from '$api/incoming/header';
	import { connected } from '$api/websocket';
</script>

<section>
	<h1>
		{#if $headerState.status === 'running'}
			<svg class="icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
				<path
					d="M16.6582 9.28638C18.098 10.1862 18.8178 10.6361 19.0647 11.2122C19.2803 11.7152 19.2803 12.2847 19.0647 12.7878C18.8178 13.3638 18.098 13.8137 16.6582 14.7136L9.896 18.94C8.29805 19.9387 7.49907 20.4381 6.83973 20.385C6.26501 20.3388 5.73818 20.0469 5.3944 19.584C5 19.053 5 18.1108 5 16.2264V7.77357C5 5.88919 5 4.94701 5.3944 4.41598C5.73818 3.9531 6.26501 3.66111 6.83973 3.6149C7.49907 3.5619 8.29805 4.06126 9.896 5.05998L16.6582 9.28638Z"
					stroke="#2ec27e"
					stroke-width="2"
					stroke-linejoin="round"
					style="--darkreader-inline-stroke: #4ed597;"
					data-darkreader-inline-stroke=""
				></path>
			</svg>
		{:else if $headerState.status === 'paused'}
			<svg class="icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
				<path
					d="M14 9V15M10 9V15M12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12C21 16.9706 16.9706 21 12 21Z"
					stroke="#e66100"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					style="--darkreader-inline-stroke: #ff842b;"
					data-darkreader-inline-stroke=""
				></path>
			</svg>
		{:else}
			<svg class="icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
				<path
					d="M12 16H12.01M12 8V12M12 21C16.9706 21 21 16.9706 21 12C21 7.02944 16.9706 3 12 3C7.02944 3 3 7.02944 3 12C3 16.9706 7.02944 21 12 21Z"
					stroke="#e01b24"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					style="--darkreader-inline-stroke: #e7353d;"
					data-darkreader-inline-stroke=""
				></path>
			</svg>
		{/if}
		<div class="title">
			{$headerState.path}
		</div>
		<div class="connection">
			{#if !$connected}
				<div class="no">Connection lost</div>
			{:else}
				<div class="yes">Connected</div>
			{/if}
		</div>
	</h1>
</section>

<style>
	.title {
		font-size: 1.5em;
		color: #69ff94;
	}
	section {
		position: sticky;
		top: 0;
		padding: 10px;
		margin: 10px;
		margin-bottom: 0;
		z-index: 100;
	}

	h1 {
		display: flex;
		align-items: center;
		gap: 10px;
		margin: 0;
		width: 100vw;
	}

	.connection {
		margin-left: auto; /* Pushes connection to the right side */
		font-size: 1em;
		width: 20%;
		color: #ff69b4; /* Optional: To make the connection status stand out */
	}
	.connection .no {
		color: #ff5555;
	}
	.connection .yes {
		color: #50fa7b;
	}

	.icon {
		width: 44px;
		height: 44px;
		display: inline-block;
	}
</style>
