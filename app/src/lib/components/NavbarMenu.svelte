<script>
  import { currentUser } from '$lib/stores/user';
</script>

<ul class="menu {$$props.class}" {...$$props.restArguments}>
  {#if $currentUser}
    <li><a href="/users">Alunos</a></li>
    <li><a href="/account">Sua Conta</a></li>

    {#if $currentUser.role === 'editor'}
      <li>
        <a href="/classes/new">Criar aula</a>
      </li>
    {/if}
  {/if}
</ul>
