<template>
  <div v-if="recipe">
    <h1>{{ recipe.title }}</h1>
    <img :src="getImageUrl(recipe.featured_image)" alt="Recipe Image" />
    <p>{{ recipe.description }}</p>
    <!-- Other recipe details -->
  </div>
</template>

<script>
export default {
  props: ['id'], // Get the recipeId from the route params
  data() {
    return {
      recipe: null,
    };
  },
  created() {
    this.fetchRecipeDetails();
  },
  methods: {
    async fetchRecipeDetails() {
      // Fetch recipe details by recipeId
      try {
        const response = await this.$apollo.query({
          query: gql`
            query GetRecipe($id: uuid!) {
              recipes_by_pk(id: $id) {
                title
                description
                featured_image
              }
            }
          `,
          variables: { id: this.id },
        });
        this.recipe = response.data.recipes_by_pk;
      } catch (error) {
        console.error('Error fetching recipe details:', error);
      }
    },
    getImageUrl(image) {
      return `/images/${image}`;
    },
  },
};
</script>
