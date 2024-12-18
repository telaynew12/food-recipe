// import gql from 'graphql-tag';

// export const ADD_RECIPE_WITH_DETAILS = gql`
//   mutation AddRecipeWithDetails(
//     $recipe: Recipes_insert_input!,
//     $images: [Images_insert_input!]!,
//     $steps: [Steps_insert_input!]!,
//     $ingredients: [Ingredients_insert_input!]!
//   ) {
//     insert_Recipes_one(object: $recipe) {
//       id
//     }
//     insert_Images(objects: $images) {
//       affected_rows
//     }
//     insert_Steps(objects: $steps) {
//       affected_rows
//     }
//     insert_Ingredients(objects: $ingredients) {
//       affected_rows
//     }
//   }
// `;