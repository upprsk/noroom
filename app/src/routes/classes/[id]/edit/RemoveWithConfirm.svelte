<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  let confirm = false;
</script>

{#if confirm}
  <button type="button" class="btn btn-warning" on:click={() => dispatch('remove')}>
    remover
  </button>
  <button type="button" class="btn" on:click={() => (confirm = false)}>cancelar</button>
{:else}
  <button type="button" class="btn btn-warning" on:click={() => (confirm = true)}> remover </button>
{/if}
