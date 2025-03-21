-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Gyms table
CREATE TABLE IF NOT EXISTS gyms (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Exercises table (global list)
CREATE TABLE IF NOT EXISTS exercises (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL  -- e.g., "Bench Press", "Squat"
);

-- Equipment Types table (global list)
CREATE TABLE IF NOT EXISTS equipment_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL  -- e.g., "Barbell", "Dumbbell"
);

-- Gym Equipment table (instances of equipment in specific gyms)
CREATE TABLE IF NOT EXISTS gym_equipment (
    id SERIAL PRIMARY KEY,
    gym_id INTEGER REFERENCES gyms(id) ON DELETE CASCADE,
    equipment_type_id INTEGER REFERENCES equipment_types(id),
    weight DECIMAL NULL,  -- Optional weight of the equipment (NULL if not applicable)
    notes VARCHAR(255) NULL,  -- Optional additional information
    UNIQUE(gym_id, equipment_type_id, weight)  -- Allow same equipment type with different weights
);

-- Workout Sessions
CREATE TABLE IF NOT EXISTS workout_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    gym_id INTEGER REFERENCES gyms(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Workout Exercises (details of each exercise in a session)
CREATE TABLE IF NOT EXISTS workout_exercises (
    id SERIAL PRIMARY KEY,
    workout_session_id INTEGER REFERENCES workout_sessions(id) ON DELETE CASCADE,
    exercise_id INTEGER REFERENCES exercises(id),
    gym_equipment_id INTEGER REFERENCES gym_equipment(id),
    weight DECIMAL NOT NULL,  -- Weight used for this specific workout
    reps INTEGER NOT NULL,
    sets INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Pantry Items
CREATE TABLE IF NOT EXISTS pantry_items (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    quantity DECIMAL NOT NULL,
    unit VARCHAR(50) NOT NULL,  -- grams, items, liters
    threshold DECIMAL NOT NULL,  -- minimum required quantity
    calories_per_unit DECIMAL NOT NULL,
    protein_per_unit DECIMAL NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name)
);

-- Meals
CREATE TABLE IF NOT EXISTS meals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Meal Ingredients (many-to-many between meals and pantry items)
CREATE TABLE IF NOT EXISTS meal_ingredients (
    meal_id INTEGER REFERENCES meals(id) ON DELETE CASCADE,
    pantry_item_id INTEGER REFERENCES pantry_items(id),
    quantity_used DECIMAL NOT NULL,
    PRIMARY KEY (meal_id, pantry_item_id)
);

-- Shopping List View (items below threshold)
CREATE OR REPLACE VIEW shopping_list AS
SELECT 
    u.id AS user_id,
    pi.name,
    (pi.threshold - pi.quantity) AS quantity_needed,
    pi.unit
FROM pantry_items pi
JOIN users u ON pi.user_id = u.id
WHERE pi.quantity < pi.threshold;