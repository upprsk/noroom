<script lang="ts">
  import BasicCard from '$lib/components/BasicCard.svelte';
  import { currentUser } from '$lib/stores/user';

  export let data;
</script>

<BasicCard>
  <svelte:fragment slot="title">C208 - L1</svelte:fragment>

  <div class="h-5"></div>

  {#if !$currentUser}
    <p>Entre com sua conta para ver o conteudo das aulas.</p>
  {/if}

  <ul class="mx-5 flex flex-col gap-3">
    {#each data.classes as cl (cl.id)}
      <li>
        <a href="classes/{cl.id}" class="btn flex justify-between">
          <span><b>{cl.title}</b></span>
          <span>{new Date(cl.updated).toLocaleString()}</span>
        </a>
      </li>
    {/each}
  </ul>
</BasicCard>
