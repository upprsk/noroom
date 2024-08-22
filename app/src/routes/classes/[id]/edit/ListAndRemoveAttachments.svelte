<script lang="ts">
  import type { ClassModel } from '$lib/models';
  import RemoveAttachment from './RemoveAttachment.svelte';

  export let klass: ClassModel;
</script>

<h4 class="text-lg">Arquivos anexados</h4>

<div class="border-l pl-5 rounded flex flex-col gap-3 max-w-md">
  {#each klass.attachments as attachment}
    <RemoveAttachment {attachment} />
  {/each}
</div>
