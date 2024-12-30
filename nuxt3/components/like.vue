<template>
  <div class="p-6 max-w-4xl mx-auto bg-white shadow-md rounded-lg">
    <div v-if="recipes.length" class="space-y-8">
      <div v-for="(recipe, index) in recipes" :key="recipe.id" class="border-b pb-6">
        <h2 class="text-2xl font-bold mb-4">{{ recipe.title }}</h2>
        <img
          v-if="recipe.featured_image"
          :src="getImageUrl(recipe.featured_image)"
          alt="Recipe Image"
          class="w-full h-500 object-cover rounded-lg mb-4"
        />
        <div class="flex justify-start gap-6 mb-4 items-center">
          <!-- Rating Section -->
          <div class="flex items-center space-x-1">
            <span class="font-semibold">Rate:</span>
            <div class="flex">
              <span
                v-for="star in 5"
                :key="star"
                :class="{
                  'text-yellow-400': recipe.userRating >= star,
                  'text-gray-400': recipe.userRating < star
                }"
                class="cursor-pointer text-2xl"
                @click="handleRating(index, star)"
              >
                ★
              </span>
            </div>
            <span class="ml-2 text-sm text-gray-600">
              ({{ recipe.avgRating.toFixed(1) }} / 5)
            </span>
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <p class="text-center">No recipes available.</p>
    </div>
  </div>
</template>
