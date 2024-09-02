<script lang="ts">
  import type { ClassModel } from '$lib/models';
  import RemoveAttachment from './RemoveAttachment.svelte';

  export let klass: ClassModel;
</script>

<h4 class="text-lg">Arquivos anexados</h4>

<div class="flex max-w-md flex-col gap-3 rounded border-l pl-5">
  {#each klass.attachments as attachment}
    <RemoveAttachment {attachment} on:remove />
  {/each}
</div>
