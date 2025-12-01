import 'package:flutter/material.dart';
import '../../env.dart';
import '../../utils/api.dart' as api;

class ExercisesPage extends StatefulWidget {
  const ExercisesPage({super.key});

  @override
  State<ExercisesPage> createState() => _ExercisesPageState();
}

class _ExercisesPageState extends State<ExercisesPage> {
  final apiUrl = Env.instance.apiExerciseUrl;
  final ScrollController _scrollController = ScrollController();
  bool _isLoading = true;
  bool _isFetchingMore = false;
  bool _hasMore = true;
  int _offset = 0;
  final int _limit = 20;
  String? _error;
  List<dynamic> _items = [];

  @override
  void initState() {
    super.initState();
    _fetchItems();
    _scrollController.addListener(_onScroll);
  }

  void _onScroll() {
    if (_scrollController.position.pixels >=
            _scrollController.position.maxScrollExtent - 200 &&
        !_isFetchingMore &&
        _hasMore &&
        !_isLoading) {
      _fetchMoreItems();
    }
  }

  Future<void> _fetchItems() async {
    try {
      final response = await api.sendProtectedGetRequest('$apiUrl/getForUser', {
        'limit': '$_limit',
        'offset': '$_offset',
      });
      setState(() {
        _items = response['exercises'];
        _isLoading = false;
        _hasMore = (_items.length == _limit);
      });
    } catch (e) {
      setState(() {
        _error = 'Failed to load items';
        _isLoading = false;
      });
    }
  }

  Future<void> _fetchMoreItems() async {
    setState(() => _isFetchingMore = true);
    try {
      _offset += _limit;
      final response = await api.sendProtectedGetRequest('$apiUrl/getForUser', {
        'limit': '$_limit',
        'offset': '$_offset',
      });
      final newItems = response['exercises'] as List<dynamic>;
      setState(() {
        _items.addAll(newItems);
        _hasMore = (newItems.length == _limit);
        _isFetchingMore = false;
      });
    } catch (e) {
      setState(() => _isFetchingMore = false);
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
      // TODO: Handle with error modal.
    }
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
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
      body: NotificationListener<ScrollNotification>(
        onNotification: (scrollInfo) {
          if (!_isFetchingMore &&
              _hasMore &&
              scrollInfo.metrics.pixels >=
                  scrollInfo.metrics.maxScrollExtent - 200) {
            _fetchMoreItems();
          }
          return false;
        },
        child: ListView.builder(
          controller: _scrollController,
          itemCount: _items.length + (_isFetchingMore ? 1 : 0),
          itemBuilder: (context, index) {
            if (index == _items.length && _isFetchingMore) {
              return const Padding(
                padding: EdgeInsets.symmetric(vertical: 16.0),
                child: Center(child: CircularProgressIndicator()),
              );
            }
            final item = _items[index];
            return ListTile(
              title: Text(item['Name'] ?? ''),
              subtitle: Text(item['TargetMuscle'] ?? ''),
            );
          },
        ),
      ),
    );
  }
}
