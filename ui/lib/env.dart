class Env {
  static final Env _instance = Env._internal();
  final String _apiUrl = "http://localhost:8080/api/v1";
  final String _apiUserUrl = "/user";
  final String _apiExerciseUrl = "/exercise";

  Env._internal();

  static Env get instance => _instance;

  String get apiUserUrl => "$_apiUrl$_apiUserUrl";
  String get apiExerciseUrl => "$_apiUrl$_apiExerciseUrl";
}
