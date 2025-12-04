import 'dart:async';

import 'package:flutter/material.dart';
import 'package:ui/pages/models/models.dart';
import '../../utils/api.dart' as api;
import '../../env.dart';

class ActiveWorkoutPage extends StatefulWidget {
  const ActiveWorkoutPage({super.key});

  @override
  State<ActiveWorkoutPage> createState() => _ActiveWorkoutPageState();
}

class _ActiveWorkoutPageState extends State<ActiveWorkoutPage> {
  List<Exercise> _exercises = [];
  final List<Workout> _workouts = [];
  Timer? _timer;
  int _secondsElapsed = 0;
  String _workoutName = 'My New Workout';

  @override
  void initState() {
    super.initState();
    _loadExercises();
    _startTimer();
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  void _startTimer() {
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      setState(() {
        _secondsElapsed++;
      });
    });
  }

  void _showEditWorkoutSessionNameModal() {
    final nameController = TextEditingController(text: _workoutName);

    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Edit Workout Name'),
          content: TextField(
            controller: nameController,
            decoration: const InputDecoration(
              labelText: 'Workout Name',
              border: OutlineInputBorder(),
            ),
            autofocus: true,
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                final newName = nameController.text.trim();
                if (newName.isNotEmpty) {
                  setState(() {
                    _workoutName = newName;
                  });
                }
                Navigator.pop(context);
              },
              child: const Text('Save'),
            ),
          ],
        );
      },
    );
  }

  String _formatTime(int seconds) {
    final minutes = seconds ~/ 60;
    final remainingSeconds = seconds % 60;
    return '${minutes.toString().padLeft(2, '0')}:${remainingSeconds.toString().padLeft(2, '0')}';
  }

  Future<void> _loadExercises() async {
    final apiUrl = Env.instance.apiExerciseUrl;
    try {
      final response = await api.sendProtectedGetRequest('$apiUrl/getForUser', {
        'limit': '500',
        'offset': '0',
      });
      setState(() {
        _exercises = (response['exercises'] as List)
            .map((e) => Exercise(id: e['id'], name: e['name']))
            .toList();
      });
    } catch (e) {
      // TODO: Handle error.
    }
  }

  void _addSet(int i) {
    setState(() {
      _workouts[i].sets.add(
        WorkoutSet(
          reps: 0,
          weight: 0,
          type: 'normal',
          order: _workouts[i].sets.length + 1,
        ),
      );
    });
  }

  Future<void> _finishWorkout() async {
    final apiUrl = Env.instance.apiWorkoutSessionUrl;
    try {
      await api.sendProtectedPostRequest('$apiUrl/create', {
        "name": _workoutName,
        "workouts": _workouts.map((workout) => workout.toJson()).toList(),
      });
    } catch (e) {
      // TODO: Handle error.
    }
    if (mounted) {
      Navigator.pop(context);
    }
  }

  void _showExerciseSelectionModal() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (BuildContext context) {
        return DraggableScrollableSheet(
          initialChildSize: 0.6,
          minChildSize: 0.3,
          maxChildSize: 0.9,
          expand: false,
          builder: (context, scrollController) {
            return Column(
              children: [
                // Modal header
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Colors.grey[100],
                    borderRadius: const BorderRadius.vertical(
                      top: Radius.circular(20),
                    ),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text(
                        'Select Exercise',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      IconButton(
                        onPressed: () => Navigator.pop(context),
                        icon: const Icon(Icons.close),
                      ),
                    ],
                  ),
                ),
                // Exercises list
                Expanded(
                  child: ListView.builder(
                    controller: scrollController,
                    itemCount: _exercises.length,
                    itemBuilder: (context, index) {
                      final exercise = _exercises[index];
                      return ListTile(
                        title: Text(exercise.name ?? 'Unnamed Exercise'),
                        onTap: () {
                          _addWorkout(exercise.id);
                          Navigator.pop(context);
                        },
                        trailing: const Icon(Icons.add_circle_outline),
                      );
                    },
                  ),
                ),
              ],
            );
          },
        );
      },
    );
  }

  void _addWorkout(int? exerciseId) {
    if (exerciseId == null) return;
    setState(() {
      _workouts.add(
        Workout(
          exerciseID: exerciseId,
          sets: [WorkoutSet(reps: 0, weight: 0, type: 'normal', order: 1)],
        ),
      );
    });
  }

  void _showEditSetModal(int workoutIndex, int setIndex) {
    final set = _workouts[workoutIndex].sets[setIndex];
    final repsController = TextEditingController(text: set.reps.toString());
    final weightController = TextEditingController(text: set.weight.toString());

    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: Text('Edit Set ${setIndex + 1}'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: repsController,
                decoration: const InputDecoration(
                  labelText: 'Reps',
                  border: OutlineInputBorder(),
                ),
                keyboardType: TextInputType.number,
              ),
              const SizedBox(height: 16),
              TextField(
                controller: weightController,
                decoration: const InputDecoration(
                  labelText: 'Weight',
                  border: OutlineInputBorder(),
                ),
                keyboardType: TextInputType.number,
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                final newReps = int.tryParse(repsController.text) ?? set.reps;
                final newWeight =
                    double.tryParse(weightController.text) ?? set.weight;

                setState(() {
                  _workouts[workoutIndex].sets[setIndex] = WorkoutSet(
                    reps: newReps,
                    weight: newWeight,
                    type: set.type,
                    order: set.order,
                  );
                });

                Navigator.pop(context);
              },
              child: const Text('Save'),
            ),
          ],
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Active Workout'),
        backgroundColor: Colors.blue[600],
        foregroundColor: Colors.white,
        actions: [
          TextButton(
            onPressed: _finishWorkout,
            child: const Text(
              'Finish',
              style: TextStyle(
                color: Colors.white,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
        ],
      ),
      body: Column(
        children: [
          // Workout header
          Container(
            padding: const EdgeInsets.all(16),
            color: Colors.grey[100],
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                GestureDetector(
                  onTap: _showEditWorkoutSessionNameModal,
                  child: Row(
                    children: [
                      Text(
                        _workoutName,
                        style: const TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(width: 8),
                      Icon(Icons.edit, size: 16, color: Colors.grey[600]),
                    ],
                  ),
                ),
                Text(
                  _formatTime(_secondsElapsed), // Live timer display
                  style: TextStyle(fontSize: 16, color: Colors.grey[600]),
                ),
              ],
            ),
          ),
          // Exercises list
          Expanded(
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: _workouts.length,
              itemBuilder: (context, index) => _buildExerciseCard(index),
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _showExerciseSelectionModal,
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildExerciseCard(int exerciseIndex) {
    final workout = _workouts[exerciseIndex];
    // TODO: Handle case where exercise is not found.
    final exerciseName = _exercises
        .firstWhere(
          (ex) => ex.id == workout.exerciseID,
          orElse: () => Exercise(id: 0, name: 'Unknown'),
        )
        .name;

    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              exerciseName ?? 'Unnamed Exercise',
              style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 12),
            ...List.generate(workout.sets.length, (setIndex) {
              final set = workout.sets[setIndex];
              return Padding(
                padding: const EdgeInsets.only(bottom: 8),
                child: Row(
                  children: [
                    Text('Set ${setIndex + 1}:'),
                    const SizedBox(width: 16),
                    Expanded(
                      child: Text('${set.reps} reps @ ${set.weight} lbs'),
                    ),
                    IconButton(
                      onPressed: () =>
                          _showEditSetModal(exerciseIndex, setIndex),
                      icon: const Icon(Icons.edit, size: 20),
                    ),
                  ],
                ),
              );
            }),
            TextButton.icon(
              onPressed: () => _addSet(exerciseIndex),
              icon: const Icon(Icons.add),
              label: const Text('Add Set'),
            ),
          ],
        ),
      ),
    );
  }
}
