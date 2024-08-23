<script>
  import { currentUser } from '$lib/stores/user';
</script>

<ul class="menu {$$props.class}" {...$$props.restArguments}>
  <li><a href="/users">Alunos</a></li>
  {#if $currentUser}
    {#if $currentUser.role === 'editor'}
      <li>
        <a href="/classes/new">Criar aula</a>
      </li>
    {/if}
    <li><a href="/account">Sua Conta</a></li>
  {/if}
</ul>
