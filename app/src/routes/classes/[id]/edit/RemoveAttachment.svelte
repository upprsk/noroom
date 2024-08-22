<script lang="ts">
  import { enhance } from '$app/forms';

  export let attachment: string;

  let confirm = false;
</script>

<form class="justify-between flex gap-2" method="POST" action="?/removefile" use:enhance>
  <input
    type="text"
    name="attachment"
    class="input input-bordered input-xs w-full max-w-xs"
    value={attachment}
    readonly
  />

  {#if confirm}
    <button class="btn btn-warning btn-xs">remover</button>
    <button type="button" class="btn btn-xs" on:click={() => (confirm = false)}> cancelar </button>
  {:else}
    <button type="button" class="btn btn-warning btn-xs" on:click={() => (confirm = true)}>
      remover
    </button>
  {/if}
</form>
