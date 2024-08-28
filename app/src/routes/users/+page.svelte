<script lang="ts">
  import BasicAvatarNoLink from '$lib/components/BasicAvatarNoLink.svelte';
  import BasicCard from '$lib/components/BasicCard.svelte';

  export let data;
</script>

<BasicCard>
  <svelte:fragment slot="title">Todos os Membros</svelte:fragment>

  <div class="overflow-x-auto">
    <table class="table">
      <!-- head -->
      <thead>
        <tr>
          <th>Nome</th>
          <th>Matrícula/Curso</th>
          <th>Posição</th>
        </tr>
      </thead>
      <tbody>
        {#each data.users as user (user.id)}
          <!-- row 1 -->
          <tr>
            <td>
              <div class="flex items-center gap-3">
                <BasicAvatarNoLink {user} />
                <div>
                  <div class="font-bold">{user.name}</div>
                  <div class="text-sm opacity-50">@{user.username}</div>
                </div>
              </div>
            </td>
            <td>
              {user.mat} - {user.curso}
            </td>
            <td>
              <div class="flex justify-between">
                <span class:font-bold={user.role === 'editor'}>{user.role}</span>

                {#if !user.verified}
                  <div class="badge badge-warning">Não validado</div>
                {/if}
              </div>
            </td>
          </tr>
        {/each}
      </tbody>
      <!-- foot -->
      <tfoot>
        <tr>
          <th>Nome</th>
          <th>Matrícula/Curso</th>
          <th>Posição</th>
        </tr>
      </tfoot>
    </table>
  </div>
</BasicCard>
