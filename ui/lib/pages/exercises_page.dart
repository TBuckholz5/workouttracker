import 'package:flutter/material.dart';
import 'dart:convert';
import '../env.dart';
import '../utils/api.dart' as api;

class ExercisesPage extends StatefulWidget {
  const ExercisesPage({super.key});

  @override
  State<ExercisesPage> createState() => _ExercisesPageState();
}

class _ExercisesPageState extends State<ExercisesPage> {
  final apiUrl = Env.instance.apiExerciseUrl;
  bool _isLoading = true;
  String? _error;
  List<dynamic> _items = [];

  @override
  void initState() {
    super.initState();
    _fetchItems();
  }

  Future<void> _fetchItems() async {
    try {
      final response = await api.sendProtectedGetRequest('$apiUrl/getForUser', {
        'limit': '50',
        'offset': '0',
      });
      setState(() {
        _items = response['exercises'];
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = 'Failed to load items';
        _isLoading = false;
      });
    }
  }

  Future<void> _addExercise(
    String name,
    String description,
    String muscleGroup,
  ) async {
    try {
      final response = await api.sendProtectedPostRequest('$apiUrl/create', {
        "Name": name,
        "Description": description,
        "TargetMuscle": muscleGroup,
      });
      setState(() {
        _items.add(response['exercise']);
      });
    } catch (e) {
      print(e);
      // TODO: Handle with error modal.
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }
    if (_error != null) {
      return Scaffold(body: Center(child: Text(_error!)));
    }
    return Scaffold(
      appBar: AppBar(
        title: const Text('Exercises'),
        automaticallyImplyLeading: false,
        leading: IconButton(
          icon: const Icon(Icons.add),
          tooltip: 'Add Exercise',
          onPressed: () {
            showDialog(
              context: context,
              builder: (context) {
                final formKey = GlobalKey<FormState>();
                String name = '';
                String description = '';
                String muscleGroup = '';
                return AlertDialog(
                  title: const Text('Add Exercise'),
                  content: Form(
                    key: formKey,
                    child: SingleChildScrollView(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          TextFormField(
                            decoration: const InputDecoration(
                              labelText: 'Exercise Name',
                            ),
                            onChanged: (value) => name = value,
                            validator: (value) => value == null || value.isEmpty
                                ? 'Enter a name'
                                : null,
                          ),
                          TextFormField(
                            decoration: const InputDecoration(
                              labelText: 'Description',
                            ),
                            onChanged: (value) => description = value,
                          ),
                          TextFormField(
                            decoration: const InputDecoration(
                              labelText: 'Target Muscle Group',
                            ),
                            onChanged: (value) => muscleGroup = value,
                          ),
                        ],
                      ),
                    ),
                  ),
                  actions: [
                    TextButton(
                      onPressed: () => Navigator.of(context).pop(),
                      child: const Text('Cancel'),
                    ),
                    ElevatedButton(
                      onPressed: () {
                        if (formKey.currentState!.validate()) {
                          _addExercise(name, description, muscleGroup);
                          Navigator.of(context).pop();
                        }
                      },
                      child: const Text('Add'),
                    ),
                  ],
                );
              },
            );
          },
        ),
      ),
      body: ListView.builder(
        itemCount: _items.length,
        itemBuilder: (context, index) {
          final item = _items[index];
          return ListTile(
            title: Text(item['Name'] ?? ''),
            subtitle: Text(item['TargetMuscle'] ?? ''),
          );
        },
      ),
    );
  }
}
