<template>
  <div>
    <header class="bg-white border-b border-gray-200">
      <div class="mx-auto max-w-7xl px-2 sm:px-4 lg:divide-y lg:divide-gray-200 lg:px-8">
        <div class="hidden sm:block">
          <div class="">
            <nav class="-mb-px flex space-x-8" aria-label="Tabs">
              <router-link
                active-class="border-indigo-500 text-indigo-600"
                class="whitespace-nowrap border-b-2 py-4 px-1 text-sm font-medium"
                to="/settings"
                exact
                ><a>Blog</a>
              </router-link>
              <router-link
                active-class="border-indigo-500 text-indigo-600"
                class="whitespace-nowrap border-b-2 py-4 px-1 text-sm font-medium"
                to="/settings/user"
                exact
                ><a>User</a></router-link
              >
            </nav>
          </div>
        </div>
      </div>
    </header>
  </div>
</template>

<script>
export default {
  name: 'NavigationBar'
}
</script>
