<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  export let attachment: string;

  let confirm = false;
</script>

<div class="flex justify-between gap-2">
  <input
    type="text"
    name="attachment"
    class="input input-xs input-bordered w-full max-w-xs"
    value={attachment}
    readonly
  />

  {#if confirm}
    <button
      type="button"
      class="btn btn-warning btn-xs"
      on:click={() => dispatch('remove', attachment)}
    >
      remover
    </button>
    <button type="button" class="btn btn-xs" on:click={() => (confirm = false)}> cancelar </button>
  {:else}
    <button type="button" class="btn btn-warning btn-xs" on:click={() => (confirm = true)}>
      remover
    </button>
  {/if}
</div>
